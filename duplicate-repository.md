The [repository](https://github.com/source-ag/assignment-software-engineering) for the assignment 
is public and Github does not allow the creation of private forks for public repositories.

The correct way of creating a private frok by duplicating the repo is documented 
[here](https://help.github.com/articles/duplicating-a-repository/).

For this assignment the commands are:

 1. Create a bare clone of the repository.
    (This is temporary and will be removed so just do it wherever.)
    ```bash
    git clone --bare git@github.com:source-ag/assignment-software-engineering.git
    ```

 2. [Create a new private repository on Github](https://help.github.com/articles/creating-a-new-repository/) and give it a good name, for example                 `source-assignment-software-engineering`.

 3. Mirror-push your bare clone to your new `source-assignment-software-engineering` repository.
    > Replace `<your_username>` with your actual Github username in the url below.
    
    ```bash
    cd assignment-software-engineering
    git push --mirror git@github.com:<your_username>/source-assignment-software-engineering.git
    ```

 4. Remove the temporary local repository you created in step 1.
    ```bash
    cd ..
    rm -rf assignment-software-engineering
    ```
    
 5. You can now clone your `source-assignment-software-engineering` repository on your machine (in my case in the `code` folder).
    ```bash
    git clone git@github.com:<your_username>/source-assignment-software-engineering.git
    ```

 6. Add the reviewers as collaborators to your new repository: DonDebonair, clickery-buttony, PruijtMeister
