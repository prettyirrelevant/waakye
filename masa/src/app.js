const express = require("express");
const helmet = require("helmet");
const xss = require("xss-clean");
const httpStatus = require("http-status");
const cors = require("cors");
const basicAuth = require("express-basic-auth");
const config = require("./config/config");
const {
  generateSpotifyAuthenticationURL,
  generateDeezerAuthenticationURL,
  handleMusicServiceAuthentication,
} = require("./utils");
const morgan = require("./middlewares/morgan");
const {
  ApiError,
  errorConverter,
  errorHandler,
} = require("./middlewares/error");
const logger = require("./config/logger");

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

app.get("/api/oauth/:platform/link", async (req, res) => {
  if (req.params.platform === "spotify") {
    return res.status(200).json({ data: generateSpotifyAuthenticationURL() });
  }

  if (req.params.platform === "deezer") {
    return res.status(200).json({ data: generateDeezerAuthenticationURL() });
  }

  return res.status(400).json({ message: "Invalid platform provided" });
});

app.post("/api/oauth/spotify", async (req, res) => {
  const { isSuccessful, statusMsg } = await handleMusicServiceAuthentication({
    successText: "spotify token saved",
    email: config.SPOTIFY_AUTH_EMAIL,
    password: config.SPOTIFY_AUTH_PASSWORD,
    serviceName: "spotify",
    submitButtonSelector: "#login-button",
    emailSelector: "#login-username",
    passwordSelector: "#login-password",
    authUrl: generateSpotifyAuthenticationURL(),
  });
  const statusCode = isSuccessful ? 200 : 500;
  return res.status(statusCode).json({ message: statusMsg });
});

app.post("/api/oauth/deezer", async (req, res) => {
  const { isSuccessful, statusMsg } = await handleMusicServiceAuthentication({
    successText: "deezer token saved",
    email: config.DEEZER_AUTH_EMAIL,
    password: config.DEEZER_AUTH_PASSWORD,
    serviceName: "deezer",
    submitButtonSelector: "#login_form_submit",
    emailSelector: "#login_mail",
    passwordSelector: "#login_password",
    authUrl: generateDeezerAuthenticationURL(),
  });
  const statusCode = isSuccessful ? 200 : 500;
  return res.status(statusCode).json({ message: statusMsg });
});

app.use((req, res, next) => {
  next(new ApiError(httpStatus.NOT_FOUND, "Not found"));
});

app.use(errorConverter);
app.use(errorHandler);

module.exports = app;
