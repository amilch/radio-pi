import logo from './logo.svg';
import './App.css';

import {MdPauseCircleFilled, MdPauseCircleOutline, MdPlayCircleFilled, MdVolumeDown, MdVolumeUp} from "react-icons/md";

import React, {useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import PlayIcon from './PlayIcon';
import PauseIcon from './PauseIcon';

const Station = ({id, title, image_url, playing, playHandler, stopHandler}) => {
  const onClick = () => {
    if(playing) {
      stopHandler();
    } else {
      playHandler();
    }
  };

  return (
    <div className="Station" onClick={onClick}>
      <div className="ImageWrapper">
        <img src={image_url} />
        <button className={playing ? "PauseButton" : "PlayButton"}>
          {playing ? <PauseIcon /> : <PlayIcon />}
        </button>
      </div>
      <div className="StationTitle">
        {title}
      </div>
    </div>
  );
}

function App() {
  const [stations, setStations] = useState([]);
  const [playing, setPlaying] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [volume, setVolume] = useState(100);

  useEffect(() => {
    fetch("http://localhost:3000/stations")
      .then(res => {
        if(res.ok) return res.json();
      })
      .then(json => {
        console.log(json);
        setStations(json);
      })
      .catch(err =>  {
        console.error(err);
        setError(err);
      })
      .finally(() => {
        setLoading(false);
      });

    fetch("http://localhost:3000/playing")
      .then(res => {
        if(res.ok) return res.json();
      })
      .then(json => {
        console.log(json);
        setPlaying(json.data);
      })
      .catch(err => {
        console.error(err);
        setError(err);
      });
  }, []);

  return (
    <div className="App">
      <header className="Header">
        <div>Radio</div>
        <AnimatePresence>
          {playing != "" && (
            <motion.div
              className="GlobalPauseButton"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            >
              <a
                onClick={() => {
                  fetch("http://localhost:3000/stop")
                    .then((res) => {
                      if (res.ok) setPlaying("");
                    })
                    .catch((err) => {
                      console.error(err);
                    });
                }}
              >
                <MdPauseCircleFilled />
              </a>
            </motion.div>
          )}
        </AnimatePresence>
        <div class="VolumeControl">
          <div class="VolumeText">{volume != 100 ? volume : volume}</div>
          <div>
            <a
              onClick={() => {
                setVolume(volume >= 10 ? volume - 10 : 0);
                fetch("http://localhost:3000/volume?set=" + volume);
              }}
            >
              <MdVolumeDown />
            </a>
          </div>

          <div>
            <a
              onClick={() => {
                setVolume(volume <= 90 ? volume + 10 : 100);
                fetch("http://localhost:3000/volume?set=" + volume);
              }}
            >
              <MdVolumeUp />
            </a>
          </div>
        </div>
      </header>
      <div className="StationsGrid">
        {stations.map((station, index) => (
          <Station
            id={station.id}
            title={station.title}
            image_url={station.image_url}
            key={index}
            playing={playing == station.id}
            playHandler={() => {
              fetch("http://localhost:3000/play?id=" + station.id)
                .then((res) => {
                  if (res.ok) setPlaying(station.id);
                })
                .catch((err) => {
                  console.error(err);
                });
            }}
            stopHandler={() => {
              fetch("http://localhost:3000/stop")
                .then((res) => {
                  if (res.ok) setPlaying("");
                })
                .catch((err) => {
                  console.error(err);
                });
            }}
          />
        ))}
      </div>
    </div>
  );
}

export default App;
