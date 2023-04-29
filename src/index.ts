import dotenv from "dotenv";
dotenv.config();

import express, { Express, NextFunction, Request, Response } from "express";
import basicAuth, { IBasicAuthedRequest } from "express-basic-auth";
import expressWinston from "express-winston";
import helmet from "helmet";
import winston from "winston";
import { envSchema, validateEnv } from "./env";
import {
  encryptWithAES256CBC,
  generateAuthenticationURL,
  handleMusicServiceAuthentication,
} from "./utils";

const app: Express = express();
const env = validateEnv(process.env, envSchema);

app.use(
  expressWinston.logger({
    transports: [new winston.transports.Console()],
    format: winston.format.combine(
      winston.format.timestamp({ format: "YYYY-MM-DD HH:mm:ss:ms" }),
      winston.format.colorize({ all: true }),
      winston.format.printf(
        (info) => `[${info.level}]: ${info.timestamp} - ${info.message}`
      )
    ),
    meta: false,
    expressFormat: true,
    colorize: true,
  })
);
app.use(helmet());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(
  "/api",
  basicAuth({
    users: { user: env.SECRET_KEY },
    unauthorizedResponse: (req: IBasicAuthedRequest) => {
      return {
        status: false,
        message: req.auth
          ? "Credentials " +
            req.auth.user +
            ":" +
            req.auth.password +
            " rejected"
          : "No credentials provided",
      };
    },
  })
);

app.get("/ping", (req: Request, res: Response) => {
  return res.status(200).json();
});

app.post("/api/oauth/spotify", async (req: Request, res: Response) => {
  const authURL = generateAuthenticationURL(
    "https://accounts.spotify.com/authorize",
    {
      response_type: "code",
      client_id: env.SPOTIFY_CLIENT_ID,
      redirect_uri: env.SPOTIFY_AUTH_REDIRECT_URI,
      scope: "playlist-modify-public",
      state: encryptWithAES256CBC(
        env.SECRET_KEY,
        env.INITIALIZATION_VECTOR,
        `${Date.now()}:spotify`
      ),
    }
  );
  const success = await handleMusicServiceAuthentication({
    successText: "spotify token saved",
    email: env.SPOTIFY_AUTH_EMAIL,
    password: env.SPOTIFY_AUTH_PASSWORD,
    serviceName: "spotify",
    submitButtonSelector: "#login-button",
    emailSelector: "#login-username",
    passwordSelector: "#login-password",
    authUrl: authURL,
  });
  return res.status(200).json({ status: success, message: null });
});

app.post("/api/oauth/deezer", async (req: Request, res: Response) => {
  const authURL = generateAuthenticationURL(
    "https://connect.deezer.com/oauth/auth.php",
    {
      app_id: env.DEEZER_APP_ID,
      redirect_uri: env.DEEZER_AUTH_REDIRECT_URI,
      perms: "manage_library,offline_access",
    }
  );
  const success = await handleMusicServiceAuthentication({
    successText: "deezer token saved",
    email: env.DEEZER_AUTH_EMAIL,
    password: env.DEEZER_AUTH_PASSWORD,
    serviceName: "deezer",
    submitButtonSelector: "#login_form_submit",
    emailSelector: "#login_mail",
    passwordSelector: "#login_password",
    authUrl: authURL,
  });
  return res.status(200).json({ status: success, message: null });
});

app.use(
  expressWinston.errorLogger({
    transports: [new winston.transports.Console()],
    format: winston.format.combine(
      winston.format.timestamp({ format: "YYYY-MM-DD HH:mm:ss:ms" }),
      winston.format.colorize({ all: true }),
      winston.format.printf(
        (info) => `[${info.level}]: ${info.timestamp} - ${info.message}`
      )
    ),
  })
);

// custom 404
app.use((req: Request, res: Response, next: NextFunction) => {
  res.status(404).json({ status: false, message: "Not Found" });
});

// custom error handler
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  res.status(500).send({ status: false, message: "Internal server error" });
});

const server = app.listen(env.PORT, () => {
  console.info(`[info]: Server is running at http://localhost:${env.PORT}`);
});

process.on("SIGTERM", (msg) => {
  console.info(`Received SIGTERM with message: ${msg}...`);
  server.close((err) => {
    console.info("Server closed.");
    process.exit(err ? 1 : 0);
  });
});
