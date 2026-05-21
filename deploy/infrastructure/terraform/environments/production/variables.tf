variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-1"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
}

variable "private_subnet_cidrs" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "public_subnet_cidrs" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

variable "general_node_desired_size" {
  description = "Desired number of general purpose nodes"
  type        = number
  default     = 5
}

variable "general_node_min_size" {
  description = "Minimum number of general purpose nodes"
  type        = number
  default     = 3
}

variable "general_node_max_size" {
  description = "Maximum number of general purpose nodes"
  type        = number
  default     = 10
}

variable "general_node_instance_types" {
  description = "Instance types for general purpose nodes"
  type        = list(string)
  default     = ["m6i.xlarge", "m6i.2xlarge"]
}

variable "memory_node_desired_size" {
  description = "Desired number of memory optimized nodes"
  type        = number
  default     = 3
}

variable "memory_node_min_size" {
  description = "Minimum number of memory optimized nodes"
  type        = number
  default     = 2
}

variable "memory_node_max_size" {
  description = "Maximum number of memory optimized nodes"
  type        = number
  default     = 8
}

variable "memory_node_instance_types" {
  description = "Instance types for memory optimized nodes"
  type        = list(string)
  default     = ["r6i.xlarge", "r6i.2xlarge"]
}

variable "burst_node_max_size" {
  description = "Maximum number of burst nodes for flash sales"
  type        = number
  default     = 20
}

variable "burst_node_instance_types" {
  description = "Instance types for burst nodes"
  type        = list(string)
  default     = ["c6i.xlarge", "c6i.2xlarge"]
}
