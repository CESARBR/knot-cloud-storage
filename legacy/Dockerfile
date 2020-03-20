FROM node:10

WORKDIR /usr/src/app

COPY package*.json ./
RUN npm i

COPY . .
RUN npm run build
RUN npm prune --production

EXPOSE 80
HEALTHCHECK CMD curl --fail http://localhost:80/healthcheck || exit 1

ENV NODE_ENV=production
CMD [ "npm", "start" ]
