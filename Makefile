deploy:
	cd infrastructure && cdk deploy --verbose --require-approval never
destroy:
	cd infrastructure && cdk destroy --force
diff:
	cd infrastructure && cdk diff