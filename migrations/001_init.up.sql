CREATE EXTENSION vector WITH SCHEMA public;

-- Creation of the Organizations table
CREATE TABLE Organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address VARCHAR(255),
    contact_information VARCHAR(255)
);

-- Creation of the Organization_Configs table
CREATE TABLE Organization_Configs (
    id SERIAL PRIMARY KEY,
    organization_id INT UNIQUE,
    web_prompt TEXT,
    apps_prompt TEXT,
    FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

-- Creation of the Users table
CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    organization_id INT,
    email VARCHAR(255),
    role VARCHAR(255),
    password VARCHAR(255), -- Added password field
    FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

-- Creation of the Chat_Session table
CREATE TABLE Chat_Session (
    id SERIAL PRIMARY KEY,
    organization_id INT,
    source VARCHAR(255), --IG, Whatsapp
    history JSON, 
    identifier VARCHAR(255), -- IG: username, Whatsapp: phone number
    ai_responder BOOLEAN,
    FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

-- Create the Indexed_Data
CREATE TABLE Indexed_Data (
  id SERIAL PRIMARY KEY,
  organization_id INT,
  data TEXT NOT NULL,
  embedding vector(1536) NOT NULL,
  FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);
