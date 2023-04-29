# Stage 1 - Build stage
FROM node:18-alpine AS build

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build


# stage 2 - Final stage
FROM node:18-slim AS final

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true

ENV BROWSER_EXECUTABLE_PATH /usr/bin/google-chrome

ENV NODE_ENV production

RUN apt-get update && apt-get install gnupg wget -y && \
    wget --quiet --output-document=- https://dl-ssl.google.com/linux/linux_signing_key.pub | gpg --dearmor > /etc/apt/trusted.gpg.d/google-archive.gpg && \
    sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' && \
    apt-get update && \
    apt-get install google-chrome-stable -y --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/app

COPY package*.json ./

RUN npm install --omit=dev

COPY --from=build /app/dist ./dist

EXPOSE 5001


CMD [ "npm", "run", "start" ]
