import React, { useState, useEffect, useRef } from 'react'
import {
  Box,
  Typography,
  TextField,
  IconButton,
  Paper,
  Avatar,
  Container,
  AppBar,
  Toolbar
} from '@mui/material'
import {
  Send as SendIcon,
  ArrowBack as ArrowBackIcon,
  People as PeopleIcon,
  EmojiEmotions as EmojiIcon,
  AttachFile as AttachIcon
} from '@mui/icons-material'
import './Room.scss'

interface Message {
  id: string
  content: string
  userId: string
  userEmail: string
  userName: string
  timestamp: Date
  isOwnMessage: boolean
}

interface RoomInfo {
  id: string
  name: string
  memberCount: number
  description?: string
}

interface RoomProps {
  id: string
  name: string
  onBack?: () => void
  currentUser?: {
    id: string
    email: string
  }
}

const Room: React.FC<RoomProps> = ({ 
  id,
  name, 
  onBack,
  currentUser = { id: "user1", email: "current@user.com" }
}) => {
  const [messages, setMessages] = useState<Message[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [roomInfo, setRoomInfo] = useState<RoomInfo | null>(null)
  const [isTyping] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  

  
  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const handleSendMessage = async () => {
    if (!newMessage.trim()) return

    const message: Message = {
      id: Date.now().toString(),
      content: newMessage.trim(),
      userId: currentUser.id,
      userEmail: currentUser.email,
      userName: currentUser.name,
      timestamp: new Date(),
      isOwnMessage: true
    }

    try {
      setMessages(prev => [...prev, message])
      setNewMessage('')
      
      
      setTimeout(() => {
        const responseMessage: Message = {
          id: (Date.now() + 1).toString(),
          content: `Thanks for your message: "${message.content}"`,
          userId: 'bot',
          userEmail: 'bot@system.com',
          userName: 'Chat Bot',
          timestamp: new Date(),
          isOwnMessage: false
        }
        setMessages(prev => [...prev, responseMessage])
      }, 1000)
      
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const formatTimestamp = (timestamp: Date) => {
    const now = new Date()
    const messageDate = new Date(timestamp)
    
    if (now.toDateString() === messageDate.toDateString()) {
      return messageDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    } else {
      return messageDate.toLocaleDateString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
    }
  }

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  const getAvatarColor = (userId: string) => {
    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F']
    const index = userId.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
    return colors[index % colors.length]
  }

  if (!roomInfo) {
    return (
      <Container className="room-loading">
        <Typography>Loading room...</Typography>
      </Container>
    )
  }

  return (
    <div className="room-container">
      <h2>{id} - {name}</h2>
      {/* Room Header */}
      <AppBar position="static" className="room-header" elevation={0}>
        <Toolbar>
          {onBack && (
            <IconButton edge="start" onClick={onBack} className="back-button">
              <ArrowBackIcon />
            </IconButton>
          )}
          <Box className="room-info">
            <Typography variant="h6" className="room-title">
              {roomInfo.name}
            </Typography>
            <Typography variant="caption" className="room-subtitle">
              <PeopleIcon fontSize="small" />
              {roomInfo.memberCount} members online
            </Typography>
          </Box>
        </Toolbar>
      </AppBar>

      {/* Messages Area */}
      <div className="messages-container">
        <div className="messages-list">
          {messages.map((message, index) => {
            const showAvatar = index === 0 || 
              messages[index - 1]?.userId !== message.userId ||
              (new Date(message.timestamp).getTime() - new Date(messages[index - 1]?.timestamp || 0).getTime()) > 300000 

            return (
              <div
                key={message.id}
                className={`message-wrapper ${message.isOwnMessage ? 'own-message' : 'other-message'}`}
              >
                <div className="message-content">
                  {!message.isOwnMessage && showAvatar && (
                    <Avatar
                      className="message-avatar"
                      sx={{ bgcolor: getAvatarColor(message.userId) }}
                    >
                      {getInitials(message.userName)}
                    </Avatar>
                  )}
                  
                  <div className="message-bubble-container">
                    {!message.isOwnMessage && showAvatar && (
                      <div className="message-header">
                        <Typography variant="caption" className="sender-name">
                          {message.userName}
                        </Typography>
                        <Typography variant="caption" className="message-time">
                          {formatTimestamp(message.timestamp)}
                        </Typography>
                      </div>
                    )}
                    
                    <Paper
                      className={`message-bubble ${message.isOwnMessage ? 'own-bubble' : 'other-bubble'}`}
                      elevation={1}
                    >
                      <Typography variant="body2" className="message-text">
                        {message.content}
                      </Typography>
                    </Paper>
                    
                    {message.isOwnMessage && (
                      <Typography variant="caption" className="message-time own-time">
                        {formatTimestamp(message.timestamp)}
                      </Typography>
                    )}
                  </div>
                </div>
              </div>
            )
          })}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Message Input */}
      <Paper className="message-input-container" elevation={3}>
        {isTyping && (
          <Box className="typing-indicator">
            <Typography variant="caption" color="text.secondary">
              Someone is typing...
            </Typography>
          </Box>
        )}
        
        <Box className="input-wrapper">
          <IconButton className="attach-button" size="small">
            <AttachIcon />
          </IconButton>
          
          <TextField
            ref={inputRef}
            fullWidth
            variant="outlined"
            placeholder="Type your message..."
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            multiline
            maxRows={4}
            className="message-input"
            InputProps={{
              endAdornment: (
                <Box className="input-actions">
                  <IconButton size="small" className="emoji-button">
                    <EmojiIcon />
                  </IconButton>
                  <IconButton
                    onClick={handleSendMessage}
                    disabled={!newMessage.trim()}
                    className="send-button"
                    size="small"
                  >
                    <SendIcon />
                  </IconButton>
                </Box>
              )
            }}
          />
        </Box>
      </Paper>
    </div>
  )
}

export default Room
