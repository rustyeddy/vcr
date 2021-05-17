
/*
This is the RedEye Intelligent Video Network (IVN). RedEye consists of
the following modules:

Camera - Each camera serves up M-JPEG from the cameras published IP
address.  The following API will be provided:

- GET /redeye/camera/<camera-id>/config
- PUT /redeye/camera/<camera-id>/config
- GET /redeye/camera/<camera-id>/status
- PUT /redeye/camera/<camera-id>/play
- PUT /redeye/camera/<camera-id>/pause

Note, the camera-id will be used for the port number when consuming an
image stream.

A-Eye - Consumes a camera stream (M-JPeg) from the prescribed source
and applies one or more image processing and/or AI filters. The
resulting M-JPEG stream is made available as it's own M-JPEG stream.

- GET /redeye/aeye/<aeye-id>/filters
- GET /redeye/aeye/<aeye-id>/filter/<filter-id>
- PUT /redeye/aeye/<aeye-id>/filter/<filter-id>
- DEL /redeye/aeye/<aeye-id>/filter/<filter-id>

Pipes:

- GET /redeye/aeye/<aeye-id>/pipes
- GET /redeye/aeye/<aeye-id>/pipe/<pipe-id>
- PUT /redeye/aeye/<aeye-id>/pipe/<pipe-id>?<filter-id>&<filter-id>&<filter-id>
- DEL /redeye/aeye/<aeye-id>/pipe/<pipe-id>

VCR:

- GET /redeye/vcr/<vcr-id>/streams
- GET /redeye/vcr/<vcr-id>/stream/<stream-id>
- PUT /redeye/vcr/<vcr-id>/stream/<stream-id>?<camera-id>&<pipe-id>
- DEL /redeye/vcr/<vcr-id>/stream/<stream-id>

The streams are a combination of the source camera and the filters (pipeline)
*/

package redeye

