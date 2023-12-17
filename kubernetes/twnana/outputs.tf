output "master_dns" {
  value = aws_instance.master.public_dns
}

output "worker_1_dns" {
  value = aws_instance.worker["worker-1"].public_dns
}

output "worker_2_dns" {
  value = aws_instance.worker["worker-2"].public_dns
}

output "master_instance_id" {
  value = aws_instance.master.id
}


output "worker_instance_ids" {
  value = [for x in aws_instance.worker : x.id]
}
