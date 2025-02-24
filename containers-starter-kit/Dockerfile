FROM public.ecr.aws/amazonlinux/amazonlinux:latest

RUN yum install go -y -q
RUN yum install procps -y -q

# Create and switch to a non root user (Unreal server builds require this)
RUN groupadd --gid 1000 serveruser \
    && useradd --uid 1000 --gid serveruser -m serveruser \
    && usermod -aG wheel serveruser
USER serveruser

# Copy over the wrapper and your game server build
COPY --chown=serveruser:serveruser ServerBuild/ ServerBuild/
COPY --chown=serveruser:serveruser SdkGoWrapper/ SdkGoWrapper/

# Copy over the wrapper script and run it as the entry point
COPY --chown=serveruser:serveruser wrapper.sh wrapper.sh

RUN chmod +x ./wrapper.sh

CMD ./wrapper.sh
