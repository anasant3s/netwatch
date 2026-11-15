
import json
from urllib.error import HTTPError, URLError
from urllib.request import urlopen

from app.core.config import settings
from app.schemas.diagnostic_schema import IncidentRequest


class NetWatchAPIError(Exception):
    def __init__(self, message: str, status_code: int = 502):
        super().__init__(message)
        self.status_code = status_code


class NetWatchService:
    def get_incident(
        self,
        device_id: int,
        incident_id: int,
    ) -> IncidentRequest:
        url = (
            f"{settings.NETWATCH_API_URL}"
            f"/api/devices/{device_id}/incidents"
        )

        try:
            with urlopen(url, timeout=10) as response:
                incidents = json.loads(
                    response.read().decode("utf-8")
                )

        except HTTPError as error:
            raise NetWatchAPIError(
                "Failed to retrieve incidents from NetWatch API.",
                502,
            ) from error

        except URLError as error:
            raise NetWatchAPIError(
                "NetWatch API is unavailable.",
                503,
            ) from error

        except (json.JSONDecodeError, UnicodeDecodeError) as error:
            raise NetWatchAPIError(
                "Invalid response from NetWatch API.",
                502,
            ) from error

        for incident in incidents:
            if incident.get("id") == incident_id:
                return IncidentRequest.model_validate(incident)

        raise NetWatchAPIError(
            "Incident not found for this device.",
            404,
        )
