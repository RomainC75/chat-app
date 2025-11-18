import React, { useState } from "react";
import {
  Box,
  Typography,
  Button,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Container,
} from "@mui/material";
import {
  Add as AddIcon,
  Login as LoginIcon,
  Close as CloseIcon,
} from "@mui/icons-material";
import "./Rooms.scss";

import { EWsMessageOut, type IwebSocketMessageOut } from "../../types/socket.type";
import { useSocket } from "../../hooks/socket.hook";
import { useSelector } from "react-redux";
import type { RootState } from "../../store/store";
import Room from "../Room/Room";
import RoomCmp from "../RoomCmp/RoomCmp";


interface RoomsProps {
  onRoomSelect?: (roomId: string) => void;
}

const Rooms: React.FC<RoomsProps> = ({ onRoomSelect }) => {
    const {sendWsMessage} = useSocket();
    const {availableRooms, publicRoom, privateRoom} = useSelector((state: RootState)=>
        state.chat
    )

  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [joinDialogOpen, setJoinDialogOpen] = useState(false);
  const [newRoomName, setNewRoomName] = useState("");
  const [newRoomDescription, setNewRoomDescription] = useState("");
  const [joinRoomCode, setJoinRoomCode] = useState("");
  const [loading, setLoading] = useState(false);


  const handleCreateRoom = async () => {
    const msg: IwebSocketMessageOut = {
        type: EWsMessageOut.createRoom,
        content: {
            name: newRoomName,
            description: newRoomDescription
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





  return (
    <>
    {
        privateRoom ? <Room id={privateRoom.id} name={privateRoom.name}/> : <div></div>
    }
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
        {availableRooms.map((room, index) => (
          <RoomCmp key={`rooomCmp-${index}`} room={room} onRoomSelect={onRoomSelect}/>
        ))}
        {privateRoom ? <RoomCmp room={privateRoom} onRoomSelect={onRoomSelect}/> : <></>}
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
    </>
  );
};

export default Rooms;
