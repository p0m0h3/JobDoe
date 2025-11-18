# JobDoe

JobDoe is a Kubernetes-based sandbox scheduling and management system designed to simplify safe execution of programs in a Kubernetes cluster. It can be used to run batch jobs and AI generated programs.

# Glossary

## Job

A single execution unit representing a program, script submitted to JobDoe for processing.

## Sandbox

An isolated environment where each Job executes securely, with limited permissions, resource quotas, and network controls.

## Session

A session is a collection of Jobs that are executed in sequence inside a Sandbox with sandbox state preseved in between Jobs.

## Runner

A component responsible for preparing, launching, and monitoring sandboxes that execute user Jobs.

## Artifact

Any file or directory produced or consumed by a Job, such as CSVs, logs, models, serialized data, or metrics.
