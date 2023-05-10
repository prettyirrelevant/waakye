const dotenv = require("dotenv");
const path = require("path");
const Joi = require("joi");

dotenv.config({ path: path.join(__dirname, "../../.env") });

const envVarsSchema = Joi.object().keys({
  PORT: Joi.number().default(5001),
  NODE_ENV: Joi.string().valid("production", "development").required(),
  DEEZER_AUTH_EMAIL: Joi.string()
    .email()
    .required()
    .description("email address to authenticate the deezer account"),
  DEEZER_AUTH_PASSWORD: Joi.string()
    .required()
    .description("password to authenticate the deezer account."),
  DEEZER_APP_ID: Joi.string()
    .required()
    .description("Deezer App ID gotten from the dashboard"),
  DEEZER_AUTH_REDIRECT_URI: Joi.string()
    .required()
    .description("The redirect URI gotten from your Deezer Dashboard"),
  SPOTIFY_AUTH_EMAIL: Joi.string()
    .email()
    .required()
    .description("email address to authenticate the spotify account"),
  SPOTIFY_AUTH_PASSWORD: Joi.string()
    .required()
    .description("password to authenticate the spotify account"),
  SPOTIFY_CLIENT_ID: Joi.string()
    .required()
    .description("Spotify Client ID gotten from the dashboard"),
  SPOTIFY_AUTH_REDIRECT_URI: Joi.string()
    .required()
    .description("The redirect URI gotten from your Spotify Dashboard"),
  BROWSER_EXECUTABLE_PATH: Joi.string()
    .default('"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"')
    .description("Path to browser executable for automation"),
  SECRET_KEY: Joi.string()
    .hex()
    .required()
    .description("Secret Key used for encryption and authentication"),
  INITIALIZATION_VECTOR: Joi.string()
    .hex()
    .required()
    .description("Initialization Vector used for the encryption"),
});
const { value: envVars, error } = envVarsSchema
  .prefs({ errors: { label: "key" } })
  .validate(process.env, { allowUnknown: true });
if (error) {
  throw new Error(`Config validation error: ${error.message}`);
}

module.exports = envVars;
