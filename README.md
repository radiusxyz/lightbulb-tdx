# Lightbulb-TDX

This repository hosts the TD-specific components of the Lightbulb project, isolating the parts that are designed to run exclusively in the TD environment.

## Running the Application

1. **Clone the repository**

    ```bash
    git clone https://github.com/radiusxyz/lightbulb-tdx.git
    cd lightbulb-tdx
    ```

2. **Set the environment variables**

    Set the environment variables in the `.env` file. Or execute the following command to copy the `.env.example` file to `.env`:

    ```bash
    make copy-env
    ```

3. **Run the gRPC server**

    ```bash
    make serve
    ```

4. **(Optional) Run the test gRPC client**

    ```bash
    make run-client
    ```
