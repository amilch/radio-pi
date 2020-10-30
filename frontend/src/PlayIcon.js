import React from "react";

const PlayIcon = ({size}) => {
    return (
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="black"
        width={size}
        height={size}
      >
        <path d="M0 0h24v24H0z" fill="none" />
        <rect height="10" width="10" x="6" y="6" fill="white" />
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 14.5v-9l6 4.5-6 4.5z" fill="currentColor"/>
      </svg>
    );
}

export default PlayIcon;