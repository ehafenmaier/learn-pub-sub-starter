FROM rabbitmq:3.13-management

# Enable the STOMP plugin
RUN rabbitmq-plugins enable rabbitmq_stomp