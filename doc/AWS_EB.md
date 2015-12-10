# Deploy to Elastic Beanstalk

## 0. Run a demo application from the web ui.

> This is the fastest way to get a service role and instance profile
> -- <cite>[Service Roles, Instance Profiles, and User Policies][1]</cite>

## 1. Install the EB client

```sh
pip install --user awsebcli
# add $HOME/.local/bin to your PATH
```

## 2. Create EB app and environment

```sh
NAME=puffin-api
REGION="eu-central-1"
PASSWORD=$(pwgen 30 1)

eb init \
  --region eu-central-1 \
  --platform "Go 1.4 (Preconfigured - Docker)" \
  $NAME

eb create $NAME \
  --cname $NAME \
  --region $REGION \
  --database \
  --database.engine postgres \
  --database.user root \
  --database.password $PASSWORD

echo $PASSWORD
```

## 3. Deploy a new version

```sh
eb use <environment>
eb deploy
```

[1]:http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/concepts-roles.html
