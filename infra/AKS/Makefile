.PHONY: init plan apply output destroy

init:
	@terraform init

plan:
	@terraform init -upgrade
	@terraform plan

apply:
	
	@terraform apply -auto-approve

output:
	@terraform output

destroy:
	@terraform destroy -auto-approve