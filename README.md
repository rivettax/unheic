# UnHEIC - HEIC/HEIF to JPEG Conversion API

A simple HTTP API service that converts HEIC/HEIF images to JPEG format using Go.

## Features

- **HEIC/HEIF to JPEG Conversion**: Convert HEIC and HEIF images to JPEG format
- **HTTP API**: Simple REST API with POST endpoint
- **Health Check**: Built-in health monitoring endpoint
- **Error Handling**: Proper error responses and logging
- **Timeout Protection**: 30-second timeout for conversion operations

## Prerequisites

- Go 1.22.3 or later

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd unheic
```

2. Install dependencies:
```bash
go mod download
```

## Running the Service

### Start the API Server

```bash
# Run the HEIC to JPEG conversion API
go run cmd/unheicd/main.go
```

The server will start on port 8080. You should see:
```
Starting server on :8080
```

### API Endpoints

- **POST `/convert`** - Convert HEIC/HEIF to JPEG
- **GET `/health`** - Health check endpoint

## Testing the API

### 1. Health Check

First, verify the service is running:

```bash
curl http://localhost:8080/health
```

Expected response: `OK`

### 2. Convert Test Image

The repository includes a pre-generated test HEIC image in the `test_heic` folder. Test the conversion:

```bash
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: image/heic" \
  --data-binary @test_heic/test.heic \
  --output test_heic/output/converted_test.jpg
```

### 3. Verify Results

Check that the conversion worked:
```bash
# Check file sizes
ls -la test_heic/*.heic test_heic/output/*.jpg

# View image information (if you have ImageMagick installed)
identify test_heic/output/converted_test.jpg
```

## API Usage Examples

### Using curl

```bash
# Basic conversion
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: image/heic" \
  --data-binary @test_heic/test.heic \
  --output test_heic/output/output.jpg

# With verbose output
curl -v -X POST http://localhost:8080/convert \
  -H "Content-Type: image/heic" \
  --data-binary @test_heic/test.heic \
  --output test_heic/output/output.jpg
```

### Using Python

```python
import requests

# Convert HEIC to JPEG
with open('test_heic/test.heic', 'rb') as f:
    response = requests.post(
        'http://localhost:8080/convert',
        data=f,
        headers={'Content-Type': 'image/heic'}
    )

if response.status_code == 200:
    with open('test_heic/output/converted.jpg', 'wb') as f:
        f.write(response.content)
    print("Conversion successful!")
else:
    print(f"Error: {response.text}")
```

### Using JavaScript/Node.js

```javascript
const fs = require('fs');
const axios = require('axios');

async function convertHeicToJpeg(inputFile, outputFile) {
    try {
        const imageData = fs.readFileSync(inputFile);
        const response = await axios.post('http://localhost:8080/convert', imageData, {
            headers: {
                'Content-Type': 'image/heic'
            },
            responseType: 'arraybuffer'
        });
        
        fs.writeFileSync(outputFile, response.data);
        console.log('Conversion successful!');
    } catch (error) {
        console.error('Error:', error.response?.data || error.message);
    }
}

convertHeicToJpeg('test_heic/test.heic', 'test_heic/output/converted.jpg');
```

## Error Handling

The API returns appropriate HTTP status codes:

- **200 OK**: Successful conversion
- **400 Bad Request**: Invalid input or conversion error
- **500 Internal Server Error**: Server error

Error responses include descriptive messages:
```json
{
    "error": "Failed to convert image: <error details>"
}
```

## Configuration

The server can be configured by modifying the following parameters in `cmd/unheicd/main.go`:

- **Port**: Change `:8080` to your preferred port
- **Timeouts**: Adjust `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`
- **Conversion Timeout**: Modify the context timeout in `handleHeifToJpeg`

## Dependencies

- `github.com/jdeng/goheif` - HEIC/HEIF decoding library
- Standard Go libraries for HTTP server and image processing

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

**Third-Party Licenses:**
- This project uses the [jdeng/goheif](https://github.com/jdeng/goheif) package, which is MIT licensed.
- The goheif package includes and statically linked code from:
  - **libde265** (LGPL-3.0)
  - **heif** (Apache-2.0)

**Implications:**
- By using this project, you inherit the obligations of the LGPL-3.0 and Apache-2.0 licenses for the portions of code provided by those libraries.
- The MIT license applies to all original code in this repository.
- The LGPL-3.0 and Apache-2.0 licenses apply to the respective third-party code included via dependencies (notably via goheif and its vendored libraries).

**Summary of goheif License Note:**
> Note: The directories 'heif' and 'libde265' contain third-party code which is licensed under their respective licenses (Apache-2.0 and LGPL-3.0) as provided in those directories.

**In summary:**
- You may use, copy, modify, and distribute this project under the terms of the MIT license.
- If you distribute binaries or modified versions, you must also comply with the LGPL-3.0 and Apache-2.0 licenses for the relevant portions of the codebase.

