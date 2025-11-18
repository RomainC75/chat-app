import {
  Box,
  Card,
  CardContent,
  Chip,
  Grid,
  Typography,
  CardActions,
  Button,
} from "@mui/material";
import {
  People as PeopleIcon,
  AccessTime as TimeIcon,
} from "@mui/icons-material";
import type { TAvailableRoom } from "../../types/chat.type";

export type TRoomCmp = {
  room: TAvailableRoom;
  onRoomSelect?: (roomId: string) => void;
};

const RoomCmp = ({ room, onRoomSelect }: TRoomCmp) => {
    console.log('-> room', room)
  const formatTime = (date: Date) => {
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / (1000 * 60));
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (minutes < 1) return "Just now";
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    return `${days}d ago`;
  };
  const handleRoomClick = (roomId: string) => {
    if (onRoomSelect) {
      onRoomSelect(roomId);
    }
  };
  return (
    <Grid size={{ xs: 12, sm: 6, md: 4 }} key={room.id}>
      <Card className="room-card" onClick={() => handleRoomClick(room.id)}>
        <CardContent>
          <Box className="room-header">
            <Typography variant="h6" component="h2" className="room-name">
              {room.name}
            </Typography>
            {room.isPrivate && (
              <Chip
                label="Private"
                size="small"
                color="secondary"
                className="private-chip"
              />
            )}
          </Box>

          {room.description && (
            <Typography
              variant="body2"
              color="text.secondary"
              className="room-description"
            >
              {room.description}
            </Typography>
          )}

          <Box className="room-stats">
            <Box className="stat-item">
              <PeopleIcon fontSize="small" />
              <Typography variant="caption">
                {room.memberCount} members
              </Typography>
            </Box>
            <Box className="stat-item">
              <TimeIcon fontSize="small" />
              <Typography variant="caption">
                {formatTime(room.lastActivity)}
              </Typography>
            </Box>
          </Box>
        </CardContent>

        <CardActions className="room-actions">
          <Button size="small" fullWidth>
            Enter Room
          </Button>
        </CardActions>
      </Card>
    </Grid>
  );
};

export default RoomCmp;
