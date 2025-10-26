import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  CardActions,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Chip,
  Container,
} from "@mui/material";
import {
  Add as AddIcon,
  Login as LoginIcon,
  People as PeopleIcon,
  AccessTime as TimeIcon,
  Close as CloseIcon,
} from "@mui/icons-material";
import "./Rooms.scss";
import useWebSocket from "react-use-websocket";
import { EWsMessageOut, type IwebSocketMessageIn, type IwebSocketMessageOut } from "../../types/socket.type";

interface Room {
  id: string;
  name: string;
  description?: string;
  memberCount: number;
  isPrivate: boolean;
  createdAt: Date;
  lastActivity: Date;
}

interface RoomsProps {
  onRoomSelect?: (roomId: string) => void;
}

const Rooms: React.FC<RoomsProps> = ({ onRoomSelect }) => {
    const token = localStorage.getItem("token")!
    const SOCKET_URL = import.meta.env.VITE_SOCKET_URL;
    console.log("->W socker url ", SOCKET_URL)

  const [rooms, setRooms] = useState<Room[]>([]);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [joinDialogOpen, setJoinDialogOpen] = useState(false);
  const [newRoomName, setNewRoomName] = useState("");
  const [newRoomDescription, setNewRoomDescription] = useState("");
  const [joinRoomCode, setJoinRoomCode] = useState("");
  const [loading, setLoading] = useState(false);

  const { sendMessage: sendWsMessage, lastMessage } =
    useWebSocket<IwebSocketMessageIn>(
      `${SOCKET_URL}/api/chat/ws?token=${token}`
    );

  useEffect(() => {
    const mockRooms: Room[] = [
      {
        id: "1",
        name: "General Chat",
        description: "General discussion for everyone",
        memberCount: 24,
        isPrivate: false,
        createdAt: new Date("2024-01-15"),
        lastActivity: new Date(),
      },
    ];
    setRooms(mockRooms);
  }, []);

  useEffect(()=>{
    console.log("-> lastM : ", lastMessage)
  }, [lastMessage])

  const handleCreateRoom = async () => {
    const msg: IwebSocketMessageOut = {
        type: EWsMessageOut.createRoom,
        content: {
            roomName: newRoomName,
            roomDescription: newRoomDescription
        }
    }
    sendWsMessage(JSON.stringify(msg))

  };

  const handleJoinRoom = async () => {
    if (!joinRoomCode.trim()) return;

    setLoading(true);
    try {
      console.log("Joining room with code:", joinRoomCode);
      setJoinDialogOpen(false);
      setJoinRoomCode("");
    } catch (error) {
      console.error("Error joining room:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleRoomClick = (roomId: string) => {
    if (onRoomSelect) {
      onRoomSelect(roomId);
    }
  };

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

  return (
    <Container maxWidth="lg" className="rooms-container">
      <Box className="rooms-header">
        <Typography variant="h4" component="h1" gutterBottom>
          Chat Rooms
        </Typography>
        <Box className="rooms-actions">
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => setCreateDialogOpen(true)}
            className="create-room-btn"
          >
            Create Room
          </Button>
          <Button
            variant="outlined"
            startIcon={<LoginIcon />}
            onClick={() => setJoinDialogOpen(true)}
            className="join-room-btn"
          >
            Join Room
          </Button>
        </Box>
      </Box>

      <Grid container spacing={3} className="rooms-grid">
        {rooms.map((room) => (
          <Grid size={{ xs: 12, sm: 6, md: 4 }} key={room.id}>
            <Card
              className="room-card"
              onClick={() => handleRoomClick(room.id)}
            >
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
        ))}
      </Grid>

      {/* Create Room Dialog */}
      <Dialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Create New Room
          <IconButton
            onClick={() => setCreateDialogOpen(false)}
            sx={{ position: "absolute", right: 8, top: 8 }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Room Name"
            fullWidth
            variant="outlined"
            value={newRoomName}
            onChange={(e) => setNewRoomName(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Description (optional)"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={newRoomDescription}
            onChange={(e) => setNewRoomDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancel</Button>
          <Button
            onClick={handleCreateRoom}
            variant="contained"
            disabled={!newRoomName.trim() || loading}
          >
            {loading ? "Creating..." : "Create Room"}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Join Room Dialog */}
      <Dialog
        open={joinDialogOpen}
        onClose={() => setJoinDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Join Room
          <IconButton
            onClick={() => setJoinDialogOpen(false)}
            sx={{ position: "absolute", right: 8, top: 8 }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Room Code or ID"
            fullWidth
            variant="outlined"
            value={joinRoomCode}
            onChange={(e) => setJoinRoomCode(e.target.value)}
            helperText="Enter the room code or ID provided by the room creator"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setJoinDialogOpen(false)}>Cancel</Button>
          <Button
            onClick={handleJoinRoom}
            variant="contained"
            disabled={!joinRoomCode.trim() || loading}
          >
            {loading ? "Joining..." : "Join Room"}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Rooms;
