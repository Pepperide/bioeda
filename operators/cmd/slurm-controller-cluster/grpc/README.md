# SLURM CONTROLLER CLUSTER - GRPC

## Server
The server code must be installed inside the slurm controller deamon in the Slurm cluster.

The slurm controller receive the script to be submitted. It save the script in a directory and run it

## Client
The client code must be instakke in the component that want to send data to the Slurm deamon in the Slurm cluster.

The client send the sbatch script to the Slurm controller deamon of the Slurm CLuster.
