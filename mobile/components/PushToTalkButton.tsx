import React, { useState } from 'react';
import { TouchableOpacity, Image, StyleSheet } from 'react-native';

interface PushToTalkButtonProps {
  onPressIn: () => void;
  onPressOut: () => void;
}

const PushToTalkButton: React.FC<PushToTalkButtonProps> = ({ onPressIn, onPressOut }) => {
  const [isPressed, setIsPressed] = useState(false);

  const handlePressIn = () => {
    setIsPressed(true);
    onPressIn();
  };

  const handlePressOut = () => {
    setIsPressed(false);
    onPressOut();
  };

  return (
    <TouchableOpacity
      style={[styles.button, isPressed && styles.buttonPressed]}
      onPressIn={handlePressIn}
      onPressOut={handlePressOut}
    >
      <Image
        source={require('../assets/sqlogo-t.png')}
        style={styles.image}
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
  image: {
    width: 80,
    height: 80,
  },
});

export default PushToTalkButton;