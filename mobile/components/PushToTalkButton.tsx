import React, { useState } from 'react';
import { TouchableOpacity, Image, StyleSheet } from 'react-native';

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
      <Image
        source={require('../assets/sqlogo-t.png')}
        style={[styles.image, disabled && styles.imageDisabled]}
        resizeMode="contain"
      />
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    position: 'absolute',
    bottom: 20,
    alignSelf: 'center',
    padding: 10,
    borderRadius: 50,
  },
  buttonPressed: {
    backgroundColor: 'rgba(255, 255, 255, 0.3)',
  },
  buttonDisabled: {
    opacity: 0.5,
  },
  image: {
    width: 80,
    height: 80,
  },
  imageDisabled: {
    opacity: 0.5,
  },
});

export default PushToTalkButton;