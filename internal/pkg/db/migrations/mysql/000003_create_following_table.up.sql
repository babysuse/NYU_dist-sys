CREATE TABLE IF NOT EXISTS Following(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Username VARCHAR (128),
    Followeename VARCHAR (128),
    FOREIGN KEY (Username) REFERENCES Users(Username),
    FOREIGN KEY (Followeename) REFERENCES Users(Username),
    PRIMARY KEY (ID)
)
