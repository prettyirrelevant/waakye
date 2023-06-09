from gevent import monkey  # isort: skip
monkey.patch_all()  # isort: skip

import os
import re
from functools import wraps

from flask import Flask, request
from flask_caching import Cache
from marshmallow import EXCLUDE, Schema, ValidationError, fields, post_dump, post_load, validate
from werkzeug.exceptions import HTTPException
from ytmusicapi import YTMusic

application = Flask(__name__)
application.config['CACHE_TYPE'] = "SimpleCache"
application.config['CACHE_DEFAULT_TIMEOUT'] = 86400

cache = Cache(application)
ytmusic = YTMusic(os.getenv("YTMUSIC_HEADERS"))


class GetPlaylistRequestSchema(Schema):
    url = fields.Url(required=True)

    @post_load
    def transform_url(self, data, **kwargs):
        match = re.match(
            pattern="^https:\/\/music\.youtube\.com\/playlist\?list=([a-zA-Z0-9-_]+)$",
            string=data["url"],
        )
        if match is None:
            raise ValidationError(
                "Invalid playlist URL. Check that it follows the format https://music.youtube.com/playlist?list=",
                field_name="url",
            )

        data["url"] = match.group(1)
        return data


class CreatePlaylistRequestSchema(Schema):
    title = fields.Str(required=True)
    description = fields.Str(load_default=None)
    track_ids = fields.List(fields.Str(required=True), required=True)
    privacy_status = fields.Str(
        validate=validate.OneOf(("PUBLIC", "PRIVATE", "UNLISTED")),
        load_default="PUBLIC",
    )


class SearchTrackRequestSchema(Schema):
    q = fields.Str(required=True)
    search_filter = fields.Str(
        data_key="filter",
        validate=validate.OneOf(("songs", "videos", "uploads")),
        load_default="songs",
    )
    scope = fields.Str(
        validate=validate.OneOf(("library", "uploads")),
        load_default=None,
    )
    limit = fields.Int(strict=True, load_default=5)
    ignore_spelling = fields.Bool(load_default=False)


class TrackResponseSchema(Schema):
    videoId = fields.Str(required=True)
    title = fields.Str(required=True)
    artists = fields.List(fields.Raw(), required=True)

    @post_dump
    def transform_artists(self, data, **kwargs):
        data["artists"] = [x["name"] for x in data["artists"]]
        data["identifier"] = data.pop("videoId")
        return data


class PlaylistResponseSchema(Schema):
    id = fields.Str(required=True)
    title = fields.Str(required=True)
    description = fields.Str(required=False)
    tracks = fields.List(fields.Nested(TrackResponseSchema(unknown=EXCLUDE)), required=True)

    @post_dump
    def transform_data(self, data, **kwargs):
        data["identifier"] = data.pop("id")
        return data


class SearchTrackResponseSchema(Schema):
    category = fields.Str(required=True)
    resultType = fields.Str(required=True)
    videoId = fields.Str(required=True)
    title = fields.Str(required=True)
    artists = fields.List(fields.Raw(required=True), required=True)

    @post_dump
    def transform_data(self, data, **kwargs):
        data["artists"] = [x["name"] for x in data["artists"]]
        data["identifier"] = data.pop("videoId")
        data["result_type"] = data.pop("resultType")

        return data


def validate_request(schema_instance):
    """Decorator to validate the request payload against a specified schema.

    This decorator validates the request payload using the provided `schema_instance`.
    If the payload fails validation, a ValidationError response is returned.

    Args:
        schema_instance (Schema): The schema instance to validate the request payload against.

    Returns:
        callable: The decorated function.
    """

    def decorator(f):
        @wraps(f)
        def decorated_function(*args, **kwargs):
            try:
                req_data = schema_instance.load(request.get_json())
            except ValidationError as e:
                return {"message": "ValidationError", "errors": e.messages}, 422

            return f(req_data, *args, **kwargs)

        return decorated_function

    return decorator


def requires_auth(f):
    """Decorator that requires authentication with a bearer token.

    This decorator checks if the request includes a valid bearer token
    for authentication. If the token is not provided or is invalid, it
    returns an AuthenticationError response.

    Args:
        f (callable): The function to be decorated.

    Returns:
        callable: The decorated function.
    """

    @wraps(f)
    def decorator(*args, **kwargs):
        if (
            request.authorization is None
            or request.authorization.type != "bearer"
            or request.authorization.token != os.getenv("SECRET_KEY")
        ):
            return {
                "message": "AuthenticationError",
                "errors": ["Invalid bearer token provided"],
            }, 401

        return f(*args, **kwargs)

    return decorator


@application.errorhandler(Exception)
def generic_errorhandler(e):
    return {"message": "InternalServerError", "errors": [str(e)]}, 500


@application.errorhandler(HTTPException)
def http_errorhandler(e: HTTPException):
    resp = e.get_response()
    return {"message": e.name, "errors": [e.description]}, resp.status_code


@application.get("/")
def index():
    return {"message": "Welcome to ytmusicapi wrapper"}


@application.post("/playlists")
@validate_request(GetPlaylistRequestSchema())
@cache.cached()
def fetch_playlist(payload):
    playlist_schema = PlaylistResponseSchema(unknown=EXCLUDE)
    result = ytmusic.get_playlist(playlistId=payload["url"], limit=None)

    return {"data": playlist_schema.dump(result)}


@application.put("/playlists")
@requires_auth
@validate_request(CreatePlaylistRequestSchema())
def create_playlist(payload):
    result = ytmusic.create_playlist(
        title=payload["title"],
        description=payload["description"],
        privacy_status=payload["privacy_status"],
        video_ids=payload["track_ids"],
    )
    if isinstance(result, dict):
        return {"message": "PlaylistCreationError", "errors": result}, 500

    return {"data": f"https://music.youtube.com/playlist?list={result}"}


@application.post("/tracks/search")
@validate_request(SearchTrackRequestSchema())
@cache.cached(timeout=43200)
def search_track(payload):
    search_schema = SearchTrackResponseSchema(unknown=EXCLUDE, many=True)
    results = ytmusic.search(
        query=payload["q"],
        filter=payload["search_filter"],
        scope=payload["scope"],
        limit=payload["limit"],
        ignore_spelling=payload["ignore_spelling"],
    )

    try:
        search_schema.load(results)
    except ValidationError as e:
        return {"message": "ValidationError", "errors": e.messages}, 422

    return {"data": search_schema.dump(results)}
