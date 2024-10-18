# Contributing

Welcome to the project! I'm really excited to have you on board, and before we dive into the guidelines, let me share the essence of why this project was created.

## Intention behind building this project

At its core, this project embodies two important aims:

1. **Start Your Open Source Journey**: It's aimed to kickstart your open-source journey. Here, you'll learn the basics of Git and get a solid grip on the [Golang](https://go.dev/) and [HTMX](https://htmx.org/) and I strongly believe that learning and building should go hand in hand.
2. **Diving into Go and HTMX**: This application will help you in your journey of understanding a client-server architecture. And I've planned much more cool stuff to add in the near future if the project hits more number of contributors.

I'd love for you to make the most of this project - it's all about learning, helping, and growing in the open-source world.

## Table of Contents

1. [Setting up the Project and Contributing](#setting-up-the-project)
2. [Code of Conduct](#code-of-conduct)
3. [Guidelines to follow](#guidelines-to-follow)


## Setting up the Project

To setup the project locally follow the steps:
1. Fork and Star the project.
2. Clone your forked repository onto your local machine.

    ```
    git clone https://github.com/<YOUR-GITHUB-USERNAME>/ArcList.git
    ```
4. Add the upstream repository

    ```
    git clone https://github.com/acmpesuecc/ArcList.git
    ```
6. Add a new branch ( THIS STEP IS OPTIONAL and you can continue to work on the main branch )

    ```
    git checkout -b "<NEW-BRANCH-NAME>"
    ```
8. Go to the repository root directory

    ```
    cd ArcList
    go run main.go
    ```
10. Naviagte to http://localhost:8080/
11. After making changes to your codebase, push the code to your forked repo
   
    ```
    git push -u origin main
    ```
    if you are not working on the main branch
    ```
    git push -u origin <BRANCH-NAME>
    ```
11. Finally have Fun 😃 and Happy Contributing !! 🥳

<a name="code-of-conduct"></a>

## Code of Conduct

In our project, we believe in creating an open and inclusive space for everyone. To ensure a respectful and positive community, follow these key guidelines:

1. **Respect Each Other**: Treat all participants kindly and respectfully.
2. **Use Inclusive Language**: Keep your language welcoming and inclusive when communicating.
3. **Accept Constructive Feedback**: Be open to constructive criticism and focus on what's best for the community.
4. **No Unacceptable Behavior**: Avoid behaviors like harassment, trolling, insults, or anything that's inappropriate in a professional setting.

We're committed to maintaining a positive and inclusive community, and your cooperation is crucial for making this a safe and enjoyable space for everyone.

<a name="setting-up-the-project"></a>

## Guidelines for Contributions

1. **Claiming an Issue**: Before you start working on an issue, make sure it's assigned to you. We do this to avoid overlapping efforts and to ensure your hard work doesn't go to waste. Please avoid raising a PR for an issue assigned to someone else.
2. **Commit Format**: When making commits, follow this format: `tag-#issue-number: <commit-message>`. The tag should be one of these: `fix`, `feat`, `docs`, `chore`, `refactor`.
3. **PR Title**: When creating a Pull Request, the title should be in this format: `tag-#issue-number: <Title-of-PR>`. Again, use the same tags: `fix`, `feat`, `docs`, `chore`, `refactor`.
4. **Selective Staging**: Make sure you stage only the necessary commits when raising a PR.
