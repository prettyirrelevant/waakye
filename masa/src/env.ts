import z, { ZodError } from "zod";

const envSchema = z.object({
  PORT: z.coerce.number(),
  DEEZER_AUTH_EMAIL: z.string(),
  DEEZER_AUTH_PASSWORD: z.string(),
  SPOTIFY_AUTH_EMAIL: z.string(),
  SPOTIFY_AUTH_PASSWORD: z.string(),
  SPOTIFY_CLIENT_ID: z.string(),
  SPOTIFY_AUTH_REDIRECT_URI: z.string(),
  SECRET_KEY: z.string(),
  DEEZER_APP_ID: z.string(),
  DEEZER_AUTH_REDIRECT_URI: z.string(),
  INITIALIZATION_VECTOR: z.string(),
  BROWSER_EXECUTABLE_PATH: z.string().optional(),
});

export function validateEnv<T extends z.ZodObject<any>>(
  env: any,
  schema: T
): z.infer<T> {
  try {
    const validatedEnv = schema.parse(env);
    return validatedEnv;
  } catch (error) {
    if (error instanceof ZodError) {
      throw new Error(`Invalid environment variables: ${error.message}`);
    } else {
      throw error;
    }
  }
}

export { envSchema };
