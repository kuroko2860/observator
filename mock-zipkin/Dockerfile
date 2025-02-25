FROM node:alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY *.js .

EXPOSE 4111

CMD [ "node", "server.js" ]