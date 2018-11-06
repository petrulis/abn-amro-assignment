# ABN AMRO Assignment

### Disclaimer

This is experimental tool to send SMS or email messages via Amazon SES or SNS. This is not yet
production ready.

## Architecture

Service infrastructure is based on Amazon Serverless Application Model therefore API
endpoints were defined in a CloudFormation template (sam.yml).
Service includes:

1. Amazon DynamoDB database as a persistence layer with configured read and write capacity at the minimum but with
ability to scale to maximum of 10 both read and write capacity units.
2. AWS Lambda Functions served as endpoints to create new message or get historical messages by
recipientIdentifier.
3. Lambda function which is scheduled to scan the database for ready to be sent messages.
4. SQS queue which contains send requests.
5. Lambda function to send message via two different distribution channels (sms and email).
6. SNS is used as distribution channel for SMS messages.
7. SES is used as distribution channel for emails.

![Diagram](/docs/architecture.png)

## Deployment Pipeline

Deployment infrastructure employs Amazon CodeSuit for building code every time it changes and
triggering CloudFormation to create and execute change sets (sam.yml)

![Diagram](/docs/pipeline.png)