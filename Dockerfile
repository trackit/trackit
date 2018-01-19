FROM node:6-slim

ARG API_URL
ENV API_URL=${API_URL}

RUN curl -sS http://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb http://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install -y supervisor git
RUN mkdir -p /var/log/supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Setup UI
WORKDIR /webui
COPY ./ ./
RUN yarn install
# Building production mode (compressing all static assets)
RUN yarn run build

# Setup exposed ports
EXPOSE 80 443

CMD ["/usr/bin/supervisord"]
