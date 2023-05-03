import crypto from "crypto";
import { Browser, Page } from "puppeteer";
import puppeteer from "puppeteer-extra";
import StealthPlugin from "puppeteer-extra-plugin-stealth";

/**
 * Encrypts a plain text using AES-256-CBC encryption with the given secret key and IV
 * @param secretKeyHex - A hex-encoded string representing the secret key to use for encryption
 * @param ivHex - A hex-encoded string representing the initialization vector to use for encryption
 * @param plainText - The plain text to encrypt
 * @returns A base64-encoded string representing the encrypted text
 */
const encryptWithAES256CBC = (
  secretKeyHex: string,
  ivHex: string,
  plainText: string
): string => {
  const key = Buffer.from(secretKeyHex, "hex");
  const iv = Buffer.from(ivHex, "hex");
  const cipher = crypto.createCipheriv("aes-256-cbc", key, iv);

  let encryptedText = cipher.update(plainText, "utf8", "hex");
  encryptedText += cipher.final("hex");

  const encryptedTextBuffer = Buffer.from(encryptedText, "hex");
  const encryptedTextBase64 = encryptedTextBuffer.toString("base64");

  return encryptedTextBase64;
};

/**
 * Sets up a new Puppeteer browser and page with anti-detection measures enabled.
 * @returns A Promise that resolves to an array containing the new page and browser instances.
 */
const getPuppeteerSetup = async (): Promise<[page: Page, browser: Browser]> => {
  puppeteer.use(StealthPlugin());

  const browserExecutablePath =
    process.env.BROWSER_EXECUTABLE_PATH ??
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome";

  const browser = await puppeteer.launch({
    headless: "new",
    args: ["--no-sandbox", "--disable-setuid-sandbox", "--disable-gpu"],
    executablePath: browserExecutablePath,
  });

  const page = await browser.newPage();

  return [page, browser];
};

interface AuthenticationParams {
  authUrl: string;
  emailSelector: string;
  passwordSelector: string;
  submitButtonSelector: string;
  successText: string;
  email: string;
  password: string;
  serviceName: "spotify" | "deezer";
}

/**
 * Handles the authentication process for a music streaming service by opening a browser window and automating the login process using Puppeteer.
 * @param {AuthenticationParams} authenticationParams - The authentication parameters for the music streaming service.
 * @returns {Promise<boolean>} Whether the authentication was successful or not.
 */
const handleMusicServiceAuthentication = async (
  authenticationParams: AuthenticationParams
): Promise<boolean> => {
  const [page, browser] = await getPuppeteerSetup();
  try {
    await page.goto(authenticationParams.authUrl);
    console.debug(`Navigated to ${authenticationParams.authUrl}...`);

    await page.type(
      authenticationParams.emailSelector,
      authenticationParams.email,
      { delay: 100 }
    );
    console.debug(`Added email...`);

    await page.type(
      authenticationParams.passwordSelector,
      authenticationParams.password,
      { delay: 100 }
    );
    console.debug(`Added password...`);

    await Promise.all([
      page.waitForNavigation(),
      page.click(authenticationParams.submitButtonSelector, { delay: 100 }),
    ]);
    console.log(`Redirected to ${authenticationParams.authUrl}`);

    const content = await page.content();
    const isSuccessful = content.includes(authenticationParams.successText, 0);
    return isSuccessful;
  } catch (error) {
    console.error(
      `Error occurred during ${authenticationParams.serviceName} OAuth: ${error}`
    );
    throw error;
  } finally {
    await browser.close();
  }
};

/**
 * Generates an authentication URL with query string parameters.
 * @param url The base URL to append the query string parameters to.
 * @param queryParams An object containing the query string parameters to be appended to the URL.
 * @returns The URL with the appended query string parameters.
 */
const generateAuthenticationURL = (
  url: string,
  queryParams: Record<string, string>
): string => {
  const searchParams = new URLSearchParams();
  for (const [key, value] of Object.entries(queryParams)) {
    searchParams.set(key, value);
  }

  return `${url}?${searchParams.toString()}`;
};

export {
  getPuppeteerSetup,
  encryptWithAES256CBC,
  generateAuthenticationURL,
  handleMusicServiceAuthentication,
};
