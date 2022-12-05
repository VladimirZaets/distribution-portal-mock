terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "dbs3" {
  bucket = "distribution-portal"

  tags = {
    Name        = "Distribution Portal"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_acl" "dbs3acl" {
  bucket = aws_s3_bucket.dbs3.id
  acl    = "private"
}

resource "aws_ecs_cluster" "app" {
  name = "app"
}

resource "aws_ecs_service" "distribution-portal" {
  name            = "distribution-portal"
  task_definition = aws_ecs_task_definition.distribution-portal.arn
  
  cluster         = aws_ecs_cluster.app.id
  launch_type     = "FARGATE"

  desired_count = 1

  load_balancer {
   target_group_arn = aws_lb_target_group.distribution-portal.arn
   container_name   = "distribution-portal"
   container_port   = "8080"
 }

  network_configuration {
   assign_public_ip = false

   security_groups = [
     aws_security_group.egress_all.id,
     aws_security_group.ingress_api.id,
   ]

   subnets = [
     aws_subnet.private_d.id,
     aws_subnet.private_e.id,
   ]
 }
}

resource "aws_cloudwatch_log_group" "distribution-portal" {
  name = "/ecs/distribution-portal"
}

resource "aws_ecs_task_definition" "distribution-portal" {
  family = "distribution-portal"

  container_definitions = <<EOF
  [
    {
      "name": "distribution-portal",
      "image": "docker.io/vzaets/distribution-portal:0.13",
      "environment": [
                {
                    "name": "DIST_PORTAL_ENV",
                    "value": "prod"
                }
            ],
      "portMappings": [
        {
          "containerPort": 8080
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-region": "us-east-1",
          "awslogs-group": "/ecs/distribution-portal",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
  EOF

  cpu = 256
  memory = 512
  requires_compatibilities = ["FARGATE"]
  execution_role_arn = aws_iam_role.distribution-portal_task_execution_role.arn
  task_role_arn = aws_iam_role.ecs_task_role.arn
  network_mode = "awsvpc"
}

resource "aws_iam_role_policy_attachment" "task_s3" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "dynamodb" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ecs_to_dynamodb.arn
}

resource "aws_iam_policy" "ecs_to_dynamodb" {
    name = "ecs_to_dynamodb"
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                "Effect": "Allow",
                "Action": ["dynamodb:*"],
                "Resource": "*"
            }            
        ]
    })
}


resource "aws_iam_role" "ecs_task_role" {
  name = "ecs-task-role"
 
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}

resource "aws_iam_role" "distribution-portal_task_execution_role" {
  name               = "distribution-portal-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json
}

resource "aws_iam_policy" "policy" {
  name        = "dp-ecs-to-s3"
  path        = "/"
  description = "Provide access from ECS distribution-portal instancess to distribution-portal s3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
    {
        "Sid": "sid1",
        "Effect": "Allow",
        "Action": [
            "s3:ListAllMyBuckets",
            "s3:ListBucket",
            "s3:HeadBucket"
        ],
        "Resource": "*"
    },
    {
        "Sid": "sid2",
        "Effect": "Allow",
        "Action": "s3:*",
        "Resource": "*"
    }
]
  })
}

data "aws_iam_policy_document" "ecs_task_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}


resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.distribution-portal_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}


resource "aws_lb_target_group" "distribution-portal" {
  name        = "distribution-portal"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.app_vpc.id

  health_check {
    enabled = true
    path    = "/health"
  }

  depends_on = [aws_alb.distribution-portal]
}

resource "aws_alb" "distribution-portal" {
  name               = "distribution-portal-lb"
  internal           = false
  load_balancer_type = "application"

  subnets = [
    aws_subnet.public_d.id,
    aws_subnet.public_e.id,
  ]

  security_groups = [
    aws_security_group.http.id,
    aws_security_group.https.id,
    aws_security_group.egress_all.id,
  ]

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_alb_listener" "distribution-portal_http" {
  load_balancer_arn = aws_alb.distribution-portal.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.distribution-portal.arn
  }
}

output "alb_url" {
  value = "http://${aws_alb.distribution-portal.dns_name}"
}