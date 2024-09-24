#!/bin/bash
S3_KEY=handler/bootstrap

# Stop the script if any command fails
set -e

validate() {
    required_envs=("S3_BUCKET_NAME" "LAMBDA_HANDLER_NAME" "AWS_ACCESS_KEY_ID" "AWS_SECRET_ACCESS_KEY" "AWS_REGION")

    for env in "${required_envs[@]}"; do
        if [ -z "${!env}" ]; then
            echo "expected environment variable '$env' cannot be empty; exiting"
            return 1;
        fi
    done 
}

echo "Validating necessary environment variables..."
if ! validate; then 
    echo "validation failed"
    exit 1
else 
    echo "validation passed"
fi

echo "Starting deployment..."
make lambda-build-rest-api

echo "Zip the binary" 
cd bin/rest-api; zip ../../bootstrap bootstrap

echo "Add binary to s3"
cd - 
aws s3 cp bootstrap.zip s3://$S3_BUCKET_NAME/$S3_KEY

aws lambda update-function-code \
  --function-name $LAMBDA_HANDLER_NAME \
  --s3-bucket $S3_BUCKET_NAME \
  --s3-key $S3_KEY

echo "removing bootstrap"
rm bootstrap.zip