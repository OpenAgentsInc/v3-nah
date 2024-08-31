import React, { useState } from "react"
import { StyleSheet, TouchableOpacity } from "react-native"
import { MaterialIcons } from "@expo/vector-icons"

interface PushToTalkButtonProps {
  onPressIn: () => void;
  onPressOut: () => void;
  disabled?: boolean;
}

const PushToTalkButton: React.FC<PushToTalkButtonProps> = ({ onPressIn, onPressOut, disabled }) => {
  const [isPressed, setIsPressed] = useState(false);

  const handlePressIn = () => {
    if (disabled) return;
    setIsPressed(true);
    onPressIn();
  };

  const handlePressOut = () => {
    if (disabled) return;
    setIsPressed(false);
    onPressOut();
  };

  return (
    <TouchableOpacity
      style={[
        styles.button,
        isPressed && styles.buttonPressed,
        disabled && styles.buttonDisabled
      ]}
      onPressIn={handlePressIn}
      onPressOut={handlePressOut}
      disabled={disabled}
    >
      <MaterialIcons
        name="mic"
        size={40}
        color={disabled ? '#666' : (isPressed ? '#ff4081' : '#fff')}
      />
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    position: 'absolute',
    bottom: 60,
    alignSelf: 'center',
    backgroundColor: '#333',
    padding: 20,
    borderRadius: 50,
  },
  buttonPressed: {
    backgroundColor: '#555',
  },
  buttonDisabled: {
    opacity: 0.5,
  },
});

export default PushToTalkButton;
