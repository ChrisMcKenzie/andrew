// package main
import * as botkit from 'botkit'
import * as apiai from 'apiai'

var ai = apiai(process.env.APIAI_TOKEN);

var controller = botkit.slackbot();

var bot = controller.spawn({
  token: process.env.SLACK_TOKEN
})

bot.startRTM(function(err,bot,payload) {
  if (err) {
    throw new Error('Could not connect to Slack');
  }
});

controller.on('direct_mention,direct_message',function(bot, message) {
  console.log(message);
  var request = ai.textRequest(message.text, {
      sessionId: message.user,
      contexts: [
        {
          name: "user",
          parameters: {
            "id": message.user,
          }
        }
      ]
  }).on('response', function(response) {
    console.log(response.result.fulfillment.messages);
    bot.reply(message, response.result.fulfillment.speech)
  }).on('error', function(error) {
      console.log(error);
  }).end();
});
