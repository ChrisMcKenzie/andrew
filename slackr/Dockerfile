FROM node

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Bundle app
COPY package.json index.ts /usr/src/app/

# Install app dependencies
RUN npm install --production --silent

CMD [ "npm", "start" ]
