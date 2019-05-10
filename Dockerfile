FROM node:12 AS BUILD_UI

RUN npm install -g @vue/cli

WORKDIR /app
COPY ui/ .
RUN yarn build

FROM nginx

RUN rm -rf /usr/share/nginx/html/*
COPY --from=BUILD_UI /app/dist/ /usr/share/nginx/html/
