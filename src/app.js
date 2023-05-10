const express = require("express");
const helmet = require("helmet");
const xss = require("xss-clean");
const httpStatus = require("http-status");
const cors = require("cors");
const basicAuth = require("express-basic-auth");
const config = require("./config/config");
const {
  encryptWithAES256CBC,
  generateAuthenticationURL,
  handleMusicServiceAuthentication,
} = require("./utils");
const morgan = require("./middlewares/morgan");
const {
  ApiError,
  errorConverter,
  errorHandler,
} = require("./middlewares/error");

const app = express();

app.use(morgan.successHandler);
app.use(morgan.errorHandler);
app.use(helmet());
app.use(
  "/api",
  basicAuth({
    users: { admin: config.SECRET_KEY },
    unauthorizedResponse: (req) =>
      req.auth
        ? "Credentials " + req.auth.user + ":" + req.auth.password + " rejected"
        : "No credentials provided",
  })
);
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(xss());
app.use(cors());
app.options("*", cors());

app.get("/ping", (req, res) => {
  return res.status(200).json();
});

app.post("/api/oauth/spotify", async (req, res) => {
  const authURL = generateAuthenticationURL(
    "https://accounts.spotify.com/authorize",
    {
      response_type: "code",
      client_id: config.SPOTIFY_CLIENT_ID,
      redirect_uri: config.SPOTIFY_AUTH_REDIRECT_URI,
      scope: "playlist-modify-public",
      state: encryptWithAES256CBC(
        config.SECRET_KEY,
        config.INITIALIZATION_VECTOR,
        `${Date.now()}:spotify`
      ),
    }
  );
  const [success, msg] = await handleMusicServiceAuthentication({
    successText: "spotify token saved",
    email: config.SPOTIFY_AUTH_EMAIL,
    password: config.SPOTIFY_AUTH_PASSWORD,
    serviceName: "spotify",
    submitButtonSelector: "#login-button",
    emailSelector: "#login-username",
    passwordSelector: "#login-password",
    authUrl: authURL,
  });
  return res.status(200).json({ status: success, message: msg });
});

app.post("/api/oauth/deezer", async (req, res) => {
  const authURL = generateAuthenticationURL(
    "https://connect.deezer.com/oauth/auth.php",
    {
      app_id: config.DEEZER_APP_ID,
      redirect_uri: config.DEEZER_AUTH_REDIRECT_URI,
      perms: "manage_library,offline_access",
    }
  );
  const [success, msg] = await handleMusicServiceAuthentication({
    successText: "deezer token saved",
    email: config.DEEZER_AUTH_EMAIL,
    password: config.DEEZER_AUTH_PASSWORD,
    serviceName: "deezer",
    submitButtonSelector: "#login_form_submit",
    emailSelector: "#login_mail",
    passwordSelector: "#login_password",
    authUrl: authURL,
  });
  return res.status(200).json({ status: success, message: msg });
});

app.use((req, res, next) => {
  next(new ApiError(httpStatus.NOT_FOUND, "Not found"));
});

app.use(errorConverter);
app.use(errorHandler);

module.exports = app;
