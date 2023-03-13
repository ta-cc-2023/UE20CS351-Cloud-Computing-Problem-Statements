# Final Project CC Lab

## Building an E-commerce Microservices Application on Cloud

## using Docker, Kubernetes, Jenkins, and Git

### Overview:

The aim of this project is to develop an e-commerce microservices application that can be deployed on
the cloud using Docker, Kubernetes, Jenkins, and Git. The application will consist of several
microservices that will be deployed as Docker containers on a Kubernetes cluster. Jenkins will be used
for continuous integration and deployment, while Git will be used for version control.

### Objectives:

- Design and implement the microservices architecture for the application.
- Create Docker containers for each microservice.
- Use Kubernetes to orchestrate the containers locally.
- Implement a Jenkins pipeline to automate the deployment process.
- Integrate Git with Jenkins to trigger the pipeline on code changes.

### Pre-requisites:

The choice of programming language depends on the students' preference and experience. However,
some of the popular languages used for developing microservices-based applications are Java, Python,
Node.js, and Go.

Before starting the development of the e-commerce application, students should have a good
understanding of the following prerequisites:

1. Web development: Students should have a good understanding of HTML, CSS, and JavaScript
    (NodeJS etc.)
2. Microservices Architecture: They should be familiar with the concepts of service discovery,
    load balancing, and API gateway.
3. Containerization: Students should have a good understanding of containerization and how it is
    used in modern application development. They should be familiar with tools like Docker and
    Kubernetes.
4. CI/CD: Students should have a good understanding of continuous integration and continuous
    deployment (CI/CD) and how it is used to automate the deployment process. They should be
    familiar with tools like Jenkins and Git.
5. Database: They should be familiar with SQL and NoSQL databases.


By having a good understanding of these prerequisites, students can create a successful e-commerce
application using a microservices-based architecture, containerization, and automation.

### Documentation:

1. Docker documentation: https://docs.docker.com/
2. Kubernetes documentation: https://kubernetes.io/docs/home/
3. Jenkins documentation: https://www.jenkins.io/doc/
4. Git documentation: https://git-scm.com/doc ; Git tutorial for beginners:
    https://www.atlassian.com/git/tutorials
5. Deploying microservices with Kubernetes: https://kubernetes.io/docs/tutorials/kubernetes-
    basics/deploy-app/deploy-intro/
6. Dockerizing a Node.js WebApp: https://nodejs.org/en/docs/guides/nodejs-docker-webapp/
7. REST API Design: https://restfulapi.net/
8. Postman documentation for API testing etc: https://learning.postman.com/docs/

### Task Break-up:

The following is a breakdown of the tasks involved in this project, along with the estimated time
required for each task:

Task 1: Design the Microservices Architecture ( 1 hour)

- Define the different microservices that will be part of the application.
- Determine the communication protocols between the microservices.
- Plan the data model and schema for the microservices.

Deliverable: Microservices architecture document.

Task 2: Develop Microservices ( 4 hours)

- Develop the different microservices using appropriate programming languages and frameworks.
- Implement REST APIs to allow communication between the microservices.
- The app should contain different modules connected to a database to store data
- For instance, a user page, product page and order page
    o User Management: This module handles the registration, authentication, and
       authorization of users. It allows users to create accounts, login, and manage their
       profiles.
    o Product Management: This module handles the management of products. It allows
       admins to add, edit, and delete products. It also allows users to view and search for
       products.
    o Order Management: This module handles the management of orders. It allows users to
       view their order history, track their orders, and manage their orders.


```
o Review Management (optional): This module handles the management of product
reviews. It allows users to view and add reviews for products.
```
Deliverable: Code for microservices.

Task 3: Containerize Microservices using Docker (2 hours)

- Write Dockerfiles for each microservice.
- Build and test Docker images for each microservice.

Deliverable: Docker images for each microservice.

Task 4: Orchestrate Microservices using Kubernetes (2 hours)

- Create Kubernetes deployment manifests for each microservice.
- Create Kubernetes services for each microservice.
- Test and validate the Kubernetes deployment.

Deliverable: Kubernetes deployment manifests and services.

Task 5: Implement Continuous Integration and Deployment using Jenkins (2 hours)

- Set up Jenkins on a server.
- Create Jenkins jobs and corresponding Jenkinsfile for building, testing, and deploying the
    microservices.
- Configure Jenkins to monitor the Git repository for changes and trigger builds and deployments
    automatically.

Deliverable: Jenkins jobs and configuration files.

Task 6: Version Control using Git (2 hours)

- Create a Git repository for the microservices code.
- Commit and push code changes to the Git repository.
- Use Git to manage different versions and branches of the code.

Deliverable: Git repository containing the microservices code.

Total Estimated Time: 1 3 hours.

### Testing Guide

The following are some examples of test cases to test the project:

1. User registration: Test that users can successfully create an account and that their information is
    stored correctly in the database.
       o Test case: Enter valid user registration details and submit the form. Verify that the user
          account is created in the database and that the user can login successfully.
2. Product catalog: Test that the product catalog displays all products and their details correctly.
    o Test case: Verify that all products are displayed on the product catalog page and that
       their details (name, description, price, image, etc.) are correct.


3. Order page: Test that users can add and remove items from their order. They must get a
    confirmation if an order has been placed. The order must be registered in the database.
4. Error handling: Test that the application handles errors and exceptions correctly.
    o Test case: Submit an invalid form (e.g. with missing required fields) and verify that the
       application displays a clear error message.

To test the application, you can use a combination of manual and automated testing techniques. Here
are some examples of how to test the application:

1. Unit testing: Write unit tests for each microservice to ensure that individual components are
    working as expected.
2. Integration testing: Test how different microservices interact with each other by running
    integration tests.
3. API testing: Test the RESTful API endpoints to verify that they return the expected data and
    handle errors correctly.
4. Load testing: Test the application's performance and scalability by simulating high traffic
    loads. To perform load testing for an e-commerce microservices application, you can use a
    load testing tool such as Apache JMeter or Gatling.
       - Set up the tool to simulate a realistic user load on the application by defining the
          number of users, their behaviour, and the duration of the test.
5. End-to-end testing: Test the entire application flow from start to finish to ensure that all
    components are working together correctly.
To perform these tests, you can use a variety of tools and frameworks, such as JUnit, Selenium,
Postman, and Apache JMeter.

### Conclusion:

This project will provide hands-on experience in building a microservices-based application using
Docker containers and deploying it on a local Kubernetes cluster. It will also provide experience in
setting up a CI/CD pipeline using Jenkins and Git, which are important skills in the cloud computing
industry.


