
if [ -z "$KeyName" ]; then
    echo "Need to set KeyName to the name of a keypair registered in your AWS account"
    exit 1
fi

if [ -z "$StackName" ]; then
    echo "Need to set StackName for cloudformation stack"
    exit 1
fi


aws cloudformation create-stack --region us-east-1 --stack-name $StackName --template-body "file://cloudformation-instances.template" --parameters ParameterKey=KeyName,ParameterValue=$KeyName
