import React, { useState } from "react"
import { StyleSheet, TouchableOpacity, View } from "react-native"
import { MaterialIcons } from "@expo/vector-icons"
import CustomText from "./CustomText"

interface PushToTalkButtonProps {
  onPressIn: () => void;
  onPressOut: () => void;
  disabled?: boolean;
  isRecording: boolean;
  isProcessing: boolean;
}

const PushToTalkButton: React.FC<PushToTalkButtonProps> = ({
  onPressIn,
  onPressOut,
  disabled,
  isRecording,
  isProcessing
}) => {
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

  const getButtonText = () => {
    if (isRecording) return "Recording";
    if (isProcessing) return "Processing";
    return "Push to talk";
  };

  return (
    <View style={styles.container}>
      <CustomText style={styles.text}>{getButtonText()}</CustomText>
      <TouchableOpacity
        style={[
          styles.button,
          isPressed && styles.buttonPressed,
          disabled && styles.buttonDisabled
        ]}
        onPressIn={handlePressIn}
        onPressOut={handlePressOut}
        disabled={disabled}
        activeOpacity={0.7}
      >
        <MaterialIcons
          name="mic"
          size={40}
          color={disabled ? '#666' : (isPressed ? '#FF0000' : '#fff')}
        />
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    position: 'absolute',
    bottom: 40,
    left: 0,
    right: 0,
  },
  text: {
    color: '#999',
    fontSize: 16,
    marginBottom: 10,
  },
  button: {
    backgroundColor: '#333',
    padding: 20,
    borderRadius: 50,
  },
  buttonPressed: {
    backgroundColor: 'rgba(51, 51, 51, 1)',
  },
  buttonDisabled: {
    opacity: 0.8,
  },
});

export default PushToTalkButton;
