# RedEye VCR

The RedEye VCR manages incoming video streams from cameras, running
the streams through processing and AI filters, then directing the
final video to the appropriate location.

## API

The following REST endpoints represent the REDEYE VCR API.

- GET /health
- GET /info
- GET /cameras
- GET /camera/<camera-id>
- GET /filters
- GET /filter/<filter-id>
- GET /aeye
- GET /aeye/<aeye-id>
- GET /streams
- GET /stores
- GET /store/<store-id>

The above API calls are all informational. The following call will
actually get a stream originated by _<camera-id>_, then run through a
specific filter <filter-id>, the result will be a processed video
stream. 

- GET /stream/<camera-id>/<filter-id>


## Tests

```bash
curl get http://localhost:8000/api/health
```

results should look like:

```json
{ 'status': true }
```

