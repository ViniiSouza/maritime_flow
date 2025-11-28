FROM node:24.11-alpine

RUN apk add git
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY src/ ./src/

EXPOSE 4000
ENV NODE_ENV=production
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json

CMD ["npm", "start"]
