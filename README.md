# Gin Template Project

This project is a template for the Gin framework, designed to quickly create Go backend projects. It includes basic authentication features such as `login` and `account creation`. The project natively supports multiple languages and has built-in system error handling with comprehensive error processing.

## Usage

### 1. Configuration

All configurations for this project start from the `.env` file. The default configuration settings are as follows:

```
config="data/config.json"
mode="development"
gin_mode="debug"
port=8080
app_name="Template"
env_mode=development
logger_lang="en"
logger_path="logs"
logger_name="{app}-{level}-{date}.log"
```
- **config**: Specifies the path to the system configuration JSON file.
- **mode**: Defines the current development mode.
- **gin_mode**: Sets the Gin mode, default is debug.
- **port**: The port on which the service listens.
- **app_name**: The name of the service. Since this is a template project, the default is `Template`.
- **logger_lang**: The default logging language. Supported languages can be found in the `locale` folder, with each JSON file - representing a specific language configuration.
- **logger_path**: The directory where logs are stored.
- **logger_name**: The log file name template. `{app}` represents the project name (`app_name`), and `{date}` represents the project startup date.

### 2. Project Structure
All startup configurations begin in the boot folder, including database initialization and Gin route creation. The entire template follows the MVC pattern to organize the code.

```bash
├─api        # Location of all API handlers
├─boot       # Project initialization and configuration files
├─common     # Utility functions
├─core       # Encapsulation of Gin framework utilities
├─dao        # Database access layer
├─data       # System configuration and other settings
├─dto        # Data Transfer Objects (DTO), including field validation
├─global     # Initialization of globally accessible variables
├─i18n       # Multi-language support and configuration definitions
├─internal   # Internal constants and methods, such as custom error types and predefined error constants
├─locale     # Multi-language configuration files
├─logs       # Directory for storing log files
├─middleware # Gin middleware implementations
├─model      # ORM models for database operations
├─router     # Global route registration
├─service    # Service layer, which isolates DAO operations from API logic and converts internal errors into custom errors
```
