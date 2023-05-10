FROM node:18-slim

WORKDIR /opt/app

COPY package*.json ./

RUN npm install --omit=dev

COPY . .

EXPOSE 5001

CMD [ "npm", "run", "start" ]
