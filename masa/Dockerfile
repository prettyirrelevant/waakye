FROM node:18-slim

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=1

WORKDIR /opt/app

COPY ./package*.json ./

RUN npm install --omit=dev

COPY . .

EXPOSE 5001

CMD [ "npm", "run", "start" ]
