
SERVICE:
• API: Get metadata for a movie
• Database: Movie metadata database
• Interacts with services: None
• Data model type: Movie metadata


• metadata/cmd --> Contains the main function for starting the service
• metadata/internal/controller -->service logic (read the movie metadata)
• metadata/internal/handler --> API handler for a service
• metadata/internal/repository  --> Logic for accessing the movie metadata database
• metadata/pkg --> data model