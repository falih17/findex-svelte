FROM node:22-alpine

WORKDIR /app

# Install dependencies first for Docker caching
COPY package.json ./
RUN npm install

COPY . .

EXPOSE 3000

ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["npm", "run", "dev"]
