import React, { useState, useCallback } from 'react';
import ChatInput from '../components/ChatInput';
import MessageList from '../components/MessageList';
import { useChatStore } from '../stores/chatStore';

export default function NativeChatTab() {
  const { messages, sendMessage, isLoading } = useChatStore();
  const [input, setInput] = useState('');

  const handleSend = useCallback(async (text: string) => {
    if (!text.trim()) return;
    await sendMessage(text);
    setInput('');
  }, [sendMessage]);

  return (
    <div className="native-chat-tab">
      <MessageList messages={messages} />
      <ChatInput
        value={input}
        onChange={setInput}
        onSend={() => handleSend(input)}
        disabled={isLoading}
      />
    </div>
  );
}
