FROM node:alpine
WORKDIR /app/frontend
COPY frontend /app/frontend
RUN mv .env.example .env
RUN yarn install
CMD yarn start