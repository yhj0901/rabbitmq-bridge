from setuptools import setup, find_packages

setup(
    name="rabbitmq-bridge-client",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "pika>=1.3.0",
    ],
    author="Your Name",
    author_email="your.email@example.com",
    description="RabbitMQ client library for messaging between Go and Python",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/username/rabbitmq-bridge",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
) 