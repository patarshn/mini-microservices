FROM node:20
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN rm -f .env
COPY .env.docker .env
EXPOSE 8081
CMD ["node", "index.js"]